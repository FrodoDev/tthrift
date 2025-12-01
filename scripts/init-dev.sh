#!/bin/bash
# scripts/init-dev.sh
echo "🚀 Thrift Development Environment"
echo "Go Version: $(go version)"
echo "Thrift Version: $(thrift --version 2>/dev/null || echo 'Not in PATH, check /usr/local/apache-thrift/bin')"
echo ""
echo "Workspace ready at /work"
exec /bin/bash