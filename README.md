# JWT Split Token

Implementing JWT Split Token approach to generate and verify JWT tokens.

## Description

Storing confidential data in a JWT Token is a bad practice! Header and payload of a JWT token is encoded (base64) and using tools such as [jwt.io](https://jwt.io) every one can see their content.
When we have to store confidential data in a JWT token, there are some approches to mitigate the risk of disclosure and data breach.
This project is based on Split Token approche. There is a good explanation of this approache in [here](https://curity.io/resources/learn/split-token-pattern/).

## Getting Started

### Dependencies

* This project is containerized and you can simply download and run it as follow.

### Download

```
# Clone this repository
$ git clone https://github.com/MoziiiKa/jwt-split-token

# Go into the repository
$ cd jwt-split-token
```

### Run

```
docker compose up -d
```
### Use (APIs)

To register a new user
```
POST http://localhost:8090/api/v1/user-management/registration HTTP/1.1
content-type: application/json
    
{
    "name": "foo",
    "username": "bar",
    "email": "foo@example.com",
    "password": "complex_password"
}
```

To login (generating a new token)
```
POST http://localhost:8090/api/v1/token-management/token HTTP/1.1
content-type: application/json
    
{
    "email": "foo@example.com",
    "password": "complex_password"
}
```

To authorized access
```
GET http://localhost:8090/api/v1/access-management/time-ir HTTP/1.1
content-type: application/text
authorization: <copy from the response of the login API and paste here>
```

## Author

Mozaffar Kazemi  

## Version History

* 0.1
    * Initial Release

## License

This project is licensed under the MIT License.