# tabi-backend

This repo is the collection of source codes for the backends for our final year project "Tabi" at HCMUS. The backends are built with microservices architecture and hosted on AWS Lambda.

## Technology

- Programming language: Golang
- Framework: Serverless
- AWS services: Lambda, VPC, RDS, EC2, EventBridge
- 3rd party APIs: Paypal Sanbox, Firebase

## Details

- Tabi-booking: The main server of our backend, contains the logic for managing bookings.
- Tabi-file: The service to connect to S3 for file storage.
- Tabi-notification: The service implementing a CRON job for push notifications with Firebase.
- Tabi-payment: The service intergrating Paypal sandbox environment for executing transactions.

## Local development

Refer to the README.md of each service on how to setup and run the server.

## Contributions

- [Nam Vu Hoai](https://github.com/namhoai1109)
- [Hieu Nguyen](https://github.com/nibtr)

## License

The project is licensed under GNU GPLv3. See the LICENSE for more details.
