#!/bin/bash

./install_binary.sh

ROOTDIR=$(cd "$(dirname "$0")" && pwd)
ORDERER_CA=${ROOTDIR}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
ORG1_CA=${ROOTDIR}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
CHANNEL_NAME="mychannel"
CHAINCODE_NAME="basic"

while [[ $# -gt 0 ]]; do
  case $1 in
    --blob-path)
      BLOB_PATH="$2"
      shift 2
      ;;
    --hash)
      HASH="$2"
      shift 2
      ;;
    --source)
      SOURCE="$2"
      shift 2
      ;;
    --peer-address)
      PEER_ADDRESS="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

if [[ -z "$BLOB_PATH" ]]; then
  echo "Error: --blob-path is required"
  exit 1
fi
if [[ -z "$HASH" ]]; then
  echo "Error: --hash is required"
  exit 1
fi
if [[ -z "$SOURCE" ]]; then
  echo "Error: --source is required"
  exit 1
fi
if [[ -z "$PEER_ADDRESS" ]]; then
  echo "Error: --peer-address is required"
  exit 1
fi

export PATH=${ROOTDIR}/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/config
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=$ORG1_CA
export CORE_PEER_MSPCONFIGPATH=${ROOTDIR}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=$PEER_ADDRESS


peer chaincode invoke \
-o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile $ORDERER_CA \
-C $CHANNEL_NAME \
-n $CHAINCODE_NAME \
--peerAddresses $CORE_PEER_ADDRESS \
--tlsRootCertFiles $ORG1_CA \
-c "{\"function\":\"CreateAsset\", \"Args\":[\"$BLOB_PATH\", \"$HASH\", \"$SOURCE\"]}"
