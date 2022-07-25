# H2S [http to smtp]

## Description

h2s is small service for sending email messages. Work with it is very simply.   
H2S can render html template with your own data and send it to recipient.   
**First:** Put your template inside `tpls` folder. _H2S support go template syntax._
**Second:** Make http post query with payload to endpoint.

## Usage

Set env variables

```
SMTP_HOST SMTP_PORT SMTP_SSL SMTP_TSL SMTP_USERNAME SMTP_PASSWORD
```

Run service ./h2s   
Send email, example with curl

```
curl localhost:8090/ --data '{"templateName":"ecommerce.html", "data":{"name":"test1", "sum":"4584", "order":"d-12323", "products":[{"name":"LG TV sdfdsfdsfsdf", "price": "777.0", "quantity":"44", "sum":"3434 руб"},{"name":"Product 2 lololo", "price":"343434", "quantity":"dsfsdfsdf", "sum":"sdfsdf"}]}, "from":"sender@rix.ru", "to":["recipent.khv@gmail.com"],"subject":"Интернет заказ в магазине санремо"}'
```

__You can test own template without send with localhost:8090/text url__


