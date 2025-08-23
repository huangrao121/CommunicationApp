# How to generate pem key pairs

`openssl ecparam -genkey -name prime256v1 -noout -out ec256-private.pem`
to generate private pem key

`openssl ec -in ec256-private.pem -pubout -out ec256-public.pem`
to generate public pem key

`openssl ecparam -genkey -name prime256v1 -noout -out mqtt-ec256-private.pem`
to generate mqtt private pem key

`openssl ec -in mqtt-ec256-private.pem -pubout -out mqtt-ec256-public.pem`
to generate mqtt public pem key