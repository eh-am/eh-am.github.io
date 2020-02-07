#!/usr/bin/env bash  
  
set -e
  
cat > /tmp/should_fail.sh <<EOF
#!/usr/bin/env bash

false
echo "Hello from inner script"
EOF

chmod +x /tmp/should_fail.sh
                                                                                                    
/tmp/should_fail.sh
