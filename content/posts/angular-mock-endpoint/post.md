



```json
"architect": {
  "serve": {
    "builder": "@angular-devkit/build-angular:dev-server",
    "options": {
      "browserTarget": "your-application-name:build",
      "proxyConfig": "src/proxy.conf.js"
    },
```
Or optionally, update `package.json` 
```
    "start": "ng serve --proxy-config src/proxy.conf.js",
```

