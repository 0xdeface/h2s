# H2S [http to smtp]

[![Go](https://github.com/0xdeface/h2s/actions/workflows/go.yml/badge.svg)](https://github.com/0xdeface/h2s/actions/workflows/go.yml)

## Description

h2s is small service for sending email messages. Work with it is very simply.   
H2S do render html template with your own data and send it to recipient.   
**First:** Put your template inside `tpls` folder. _H2S support go template syntax._  
**Second:** Make http post query with payload to endpoint.

### Required env variables

```
SMTP_HOST 
SMTP_PORT 
SMTP_SSL 
SMTP_TSL 
SMTP_USERNAME 
SMTP_PASSWORD
```

### Required payload fields

```
{
"templateName": "target_template.html",
"data" : "data for you template",
"subject" : "email subject",
"to" :["array of recipients",]
"from" : "email from"
}
```

### Endpoints

/send [POST] - synchronous rendering html and sending to recipient   
/send-async [POST] - asynchronous returning uuid of job  
/result?uuid=<YOUR UUID OF REQUEST> [GET] - returning result previously added job /test - returning rendered template
without sending smtp message

## Usage

Run service ./h2s   
Send request. Example with curl bellow

```
curl localhost:8090/send --data '{"templateName":"ecommerce.html", "data":{"name":"test1", "sum":"4584", "order":"d-12323", "products":[{"name":"LG TV sdfdsfdsfsdf", "price": "777.0", "quantity":"44", "sum":"3434 руб"},{"name":"Product 2 lololo", "price":"343434", "quantity":"dsfsdfsdf", "sum":"sdfsdf"}]}, "from":"sender@rix.ru", "to":["recipent.khv@gmail.com"],"subject":"Интернет заказ в магазине"}'
```




