
+++
date = "2022-05-05"
title = "ALB Controller with AWS Cognito and Google Oauth"
categories = ["blog", "kubernetes", "aws"]
+++

At $WORK we needed to expose some apps that don't have a login flow (eg alertmanager, prometheus), these are running in kubernetes.

Since I found the documentation a little bit confusing I decided to create this blogpost to document to others (here myself from the future included).


Most examples here will be in terraform, which can be easily translated to the UI.
# Requirements
* ExternalDNS setup
* A domain you already know
* All required permissions to create all necessary stuff



First let's create a user pool and a domain. \
**Pay attention to the domain name, since it will be used later on.**


## 1. user pool + domain
```terraform
resource "aws_cognito_user_pool" "myproject" {
  name = "myproject"
}

resource "aws_cognito_user_pool_domain" "myproject" {
  domain       = "myproject"
  user_pool_id = aws_cognito_user_pool.myproject.id
}
```

## 2. google client
Then at google UI let's [create a Client ID](https://console.cloud.google.com/apis/credentials)

There at  **Authorized JavaScript origins** you specify your domain name you created previously
```
https://myproject.auth.us-east-1.amazoncognito.com
```

And for **Authorized redirect URIs** you specify the same thing plus the callback url
```
https://myproject.auth.us-east-1.amazoncognito.com/oauth2/idpresponse
```

Copy the `client id` and `client secret`.

<!--
## TODO criar a parada de Consent Screen
https://aws.amazon.com/premiumsupport/knowledge-center/cognito-google-social-identity-provider/#Configure_the_OAuth_consent_screen
-->


## 3. Identity provider
Now you create the google provider, the client id and secret is the one you copied in a previous step.
```terraform
resource "aws_cognito_identity_provider" "google_provider" {
  user_pool_id  = aws_cognito_user_pool.myproject.id
  provider_name = "Google"
  provider_type = "Google"

  provider_details = {
    authorize_scopes = "profile email openid"
    client_id        = "YOUR_CLIENT_ID"
    client_secret    = "YOUR_CLIENT_SECRET"
  }

  attribute_mapping = {
    email    = "email"
  }
}
```

### ignore changes
As you login state will be updated, which terraform is not aware about. To make terraform ignore it use:

```terraform
resource "aws_cognito_identity_provider" "google_provider" {
  (...)

  lifecycle {
    ignore_changes = [
      attribute_mapping["username"],
      provider_details["attributes_url"],
      provider_details["attributes_url_add_attributes"],
      provider_details["oidc_issuer"],
      provider_details["token_request_method"],
      provider_details["token_url"],
      provider_details["authorize_url"],
    ]
  }
}
```

## 3. aws cognito client
Now that we have the google provider we can create a client. Pay attention to the callback URL, since you need to append `/oauth2/idpresponse`.

```terraform
resource "aws_cognito_user_pool_client" "test_client" {
  name = "google_client"

  user_pool_id = aws_cognito_user_pool.myproject.id
  callback_urls                        = ["https://MY_EXPOSED_WEBSITE/oauth2/idpresponse"]
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = ["email", "openid"]
  supported_identity_providers         = ["Google"]
  generate_secret = true
}

```

## 4. The kubernetes bit
Let's create a simple nginx just to test this. This part should be straightforward (if you know k8s).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  type: ClusterIP
  ports:
    - port: 80
      targetPort: 80
```

Then the ingress part
```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: "nginx-alb-ingress"
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: "ip"
    alb.ingress.kubernetes.io/subnets: YOUR_SUBNETS
    alb.ingress.kubernetes.io/listen-ports: "[{\"HTTP\": 80}, {\"HTTPS\": 443}]"
    alb.ingress.kubernetes.io/ssl-redirect: "443"
    alb.ingress.kubernetes.io/certificate-arn: YOUR_CERTIFICATE
    external-dns.alpha.kubernetes.io/hostname: YOUR_HOSTNAME

    # Cognito stuff
    alb.ingress.kubernetes.io/auth-type: cognito
    alb.ingress.kubernetes.io/auth-scope: openid
    alb.ingress.kubernetes.io/auth-session-timeout: '3600'
    alb.ingress.kubernetes.io/auth-session-cookie: AWSELBAuthSessionCookie
    alb.ingress.kubernetes.io/auth-on-unauthenticated-request: authenticate
    alb.ingress.kubernetes.io/auth-idp-cognito: '{"UserPoolArn":"YOUR_POOL", "UserPoolClientId":"YOUR_CLIENT_ID", "UserPoolDomain":"YOUR_DOMAIN"}'
spec:
  rules:
  - host: YOUR_HOSTNAME
    http:
      paths:
        - path: /*
          backend:
            serviceName: "nginx-service"
            servicePort: 80
```


I think here everything should be straightforward, except the `auth-idp-cognito` part
```
alb.ingress.kubernetes.io/auth-idp-cognito: '{"UserPoolArn":"YOUR_POOL", "UserPoolClientId":"YOUR_CLIENT_ID", "UserPoolDomain":"YOUR_DOMAIN"}'
```

It's a JSON, so let's break it down

| | |
|-|-|
| `"UserPoolArn":"YOUR_POOL"` | you can get that via the UI by accessing General Settings |
`"UserPoolClientId":"YOUR_CLIENT_ID"` | same thing, but under App Integration -> App Client Settings
`"UserPoolDomain":"YOUR_DOMAIN"` | thing here is that we only want the prefix, which is what we set in `aws_cognito_user_pool_domain.myproject.domain`


And that should be it.

