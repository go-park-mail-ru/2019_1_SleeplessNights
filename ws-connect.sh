#!/bin/bash

xterm -hold -e "gws client -url "127.0.0.1:"$1"/ws/connect" -header "Cookie:session_token=$2"" &
