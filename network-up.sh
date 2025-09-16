#!/bin/bash


ROOTDIR=$(cd "$(dirname "$0")" && pwd)
export PATH=${ROOTDIR}/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx
export VERBOSE=false

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${ROOTDIR}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=${ROOTDIR}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
export PEER0_ORG2_CA=${ROOTDIR}/organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
export PEER0_ORG3_CA=${ROOTDIR}/organizations/peerOrganizations/org3.example.com/tlsca/tlsca.org3.example.com-cert.pem

CHANNEL_NAME="mychannel"
DELAY="3"
MAX_RETRY="5"
BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"

export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key


function networkUp() {
  if [ ! -d "organizations/peerOrganizations" ]; then
    rm -Rf organizations/peerOrganizations && rm -Rf organizations/ordererOrganizations

    cryptogen generate --config=./organizations/cryptogen/crypto-config-org1.yaml --output="organizations"
    cryptogen generate --config=./organizations/cryptogen/crypto-config-org2.yaml --output="organizations"
    cryptogen generate --config=./organizations/cryptogen/crypto-config-orderer.yaml --output="organizations"
  fi

  docker-compose -f compose/compose-test-net.yaml up
}

function setGlobals() {
  local USING_ORG=$1
  if [ $USING_ORG -eq 1 ]; then
    export CORE_PEER_LOCALMSPID=Org1MSP
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=${ROOTDIR}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
  elif [ $USING_ORG -eq 2 ]; then
    export CORE_PEER_LOCALMSPID=Org2MSP
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=${ROOTDIR}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
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

  # Join org1 peer to the channel 
  setGlobals 1
  peer channel join -b $BLOCKFILE
  
  # Join org2peer to the channel 
  setGlobals 2
  peer channel join -b $BLOCKFILE
}

function deployCC() {
  # Vendoring Go dependencies
  pushd ./chaincode/go
  GO111MODULE=on go mod vendor
  popd
  
   
}

#networkUp
createChannel