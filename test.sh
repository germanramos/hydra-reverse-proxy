#!/bin/sh -e

. ./build

ginkgo -r --failOnPending --cover --trace reverse_proxy/