#!/bin/bash

./install_binary.sh

SOCK="${DOCKER_HOST:-/var/run/docker.sock}"
DOCKER_SOCK="${SOCK##unix://}"

ROOTDIR=$(cd "$(dirname "$0")" && pwd)
export PATH=${ROOTDIR}/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/config
export VERBOSE=false

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${ROOTDIR}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
export ORG1_CA=${ROOTDIR}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem

CHANNEL_NAME="mychannel"
DELAY="3"
MAX_RETRY="5"
BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"
CHAINCODE_NAME="basic"

export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key

function networkUp() {
  # cleanup the old instances
  DOCKER_SOCK="${DOCKER_SOCK}" docker-compose -f compose-test-net.yaml down --volumes --remove-orphans
  docker rm -f $(docker ps -aq --filter label=service=hyperledger-fabric) 2>/dev/null || true
  docker rm -f $(docker ps -aq --filter name='dev-peer*') 2>/dev/null || true
  docker kill "$(docker ps -q --filter name=ccaas)" 2>/dev/null || true
  docker image rm -f $(docker images -aq --filter reference='dev-peer*') 2>/dev/null || true

  if [ ! -d "organizations/peerOrganizations" ]; then
    rm -Rf organizations/peerOrganizations && rm -Rf organizations/ordererOrganizations

    cryptogen generate --config=./organizations/cryptogen/crypto-config-org1.yaml --output="organizations"
    cryptogen generate --config=./organizations/cryptogen/crypto-config-orderer.yaml --output="organizations"
  fi

  DOCKER_SOCK="${DOCKER_SOCK}" docker-compose -f compose-test-net.yaml up -d
}

function setGlobals() {
  local USING_PEER=$1

  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID=Org1MSP
  export CORE_PEER_TLS_ROOTCERT_FILE=$ORG1_CA
  export CORE_PEER_MSPCONFIGPATH=${ROOTDIR}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
  
  if [ $USING_PEER -eq 0 ]; then
    export CORE_PEER_ADDRESS=localhost:7051
  elif [ $USING_PEER -eq 1 ]; then
    export CORE_PEER_ADDRESS=localhost:8051
  elif [ $USING_PEER -eq 2 ]; then
    export CORE_PEER_ADDRESS=localhost:10051
  fi
}

function createChannel() {
  if [ ! -d "channel-artifacts" ]; then
    mkdir channel-artifacts
  fi

  # Creates channel genesis block
  configtxgen -profile ChannelUsingRaft -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME

  # Create channel
  osnadmin channel join --channelID ${CHANNEL_NAME} --config-block ./channel-artifacts/${CHANNEL_NAME}.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"

  # Join peer 0 to the channel
  setGlobals 0
  peer channel join -b $BLOCKFILE

  # Join peer 1 to the channel
  setGlobals 1
  peer channel join -b $BLOCKFILE

  # Join peer 2 to the channel
  setGlobals 2
  peer channel join -b $BLOCKFILE
}

function deployCC() {
  # Vendoring Go dependencies
  pushd ./chaincode-go
  GO111MODULE=on go mod vendor
  popd

  # package chaincode
  peer lifecycle chaincode package chaincode.tar.gz --path ./chaincode-go --lang golang --label chaincode
  
  PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid chaincode.tar.gz)

  # install chaincode on peer0.org1
  setGlobals 0
  peer lifecycle chaincode install chaincode.tar.gz
  
  # install chaincode on peer1.org1
  setGlobals 1
  peer lifecycle chaincode install chaincode.tar.gz

  # install chaincode on peer2.org1
  setGlobals 2
  peer lifecycle chaincode install chaincode.tar.gz

  
  # approve chaincode
  setGlobals 1 0
  peer lifecycle chaincode approveformyorg \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls \
    --cafile "$ORDERER_CA" \
    --channelID $CHANNEL_NAME \
    --name $CHAINCODE_NAME \
    --version 1 \
    --package-id $PACKAGE_ID \
    --sequence 1

  # commit chaincode
  peer lifecycle chaincode commit \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls \
    --cafile "$ORDERER_CA" \
    --channelID $CHANNEL_NAME \
    --name $CHAINCODE_NAME \
    --version 1 \
    --sequence 1 \
    --peerAddresses localhost:7051 \
    --tlsRootCertFiles $ORG1_CA
}

function writeChaincode() {
  ./add_entry.sh --blob-path "blobPath" --hash "hash" --source "source" --peer-address "localhost:7051" 

  sleep 5
  ./read_entries.sh --source "source" --peer-address "localhost:7051"
}

networkUp
sleep 3

createChannel
deployCC