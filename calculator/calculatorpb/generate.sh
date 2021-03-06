#!/bin/bash

#Code generate script for calculator RPC
protoc ./calculator.proto --go_out=plugins=grpc:.