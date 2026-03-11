#!/bin/bash

curl -X POST -H "Content-Type: application/json" -d @questions.json http://localhost:9090/questions