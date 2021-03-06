#!/bin/bash
#Code generate script for greet RPC
protoc ./greet.proto --go_out=plugins=grpc:.