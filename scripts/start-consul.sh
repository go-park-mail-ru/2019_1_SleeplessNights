#!/bin/bash

consul agent -dev -config-dir=../consul.d -node=machine -ui -http-port 8010 -client 0.0.0.0
