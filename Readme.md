## Backoffice

DarchLabs backoffice api

### API spec


#### **GET /api/v1/health**

**Response**

Status: 200

```json
{
	"status": "running"
}
```

#### **POST /api/v1/users/signup**

**Request**

```json
{
	"email": "jdoe@gmail.com",
	"nickname": "schecoperez",
	"password": "fuck.max1234"
}
```

**Response**

Status: 201

```json
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Impkb2VAZ21haWwuY29tIiwiZXhwIjoxNzE5MTU3MjczfQ.MfaT_5lrapX4sapI992uQodW0xHsbv4UeNf0guCUEaA"
}
```

#### **POST /api/v1/users/login**

**Request**

```json
{
	"email": "jdoe@gmail.com",
	"password": "password124"
}
```

**Response**

Status: 201

```json
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Impkb2VAZ21haWwuY29tIiwiZXhwIjoxNzE5MTU3MjczfQ.MfaT_5lrapX4sapI992uQodW0xHsbv4UeNf0guCUEaA"
}
```


#### **PUT /api/v1/users/profiles/{short_id}**

Authorization: `Yes`

This endpoint expects an Authorization header with Bearer tooken schema

This endpoint works as an upsert in background. To perform an upsert the `short_id` must be present as path param. This because the `short_id` will be on the nfc card.

The server wil only update the fields present in the request.

**Request**

```json
{
	"linkedin": "https://linkedin.com/@/schecoperez",
	"email": "schp@redbullracing.com",
	"whatsapp": "+56989040598",
	"medium": "https://blog.medium.com/me/schecoperez",
	"twitterX": "https://x.com/profiles/schecoperez",
	"website": "https://elropeselaco.me"
}
```

**Response**

Status: 201

```json
{
	"profile": {
		"UserID": "3adaeaa9-9f72-421f-981f-1ce06d3561aa",
		"ShortID": "abdExsk",
		"Linkedin": "https://linkedin.com/@/schecoperez",
		"Email": "schp@redbullracing.com",
		"Whatsapp": null,
		"Medium": null,
		"TwitterX": null,
		"Website": null,
		"Image": "",
		"CreatedAt": "2024-02-27T02:57:29.600586Z",
		"UpdatedAt": "2024-02-27T02:58:21.879074Z",
		"DeletedAt": null
	},
	"status": "updated"
}
```


#### **Get /api/v1/users/profiles**

This endpoint enables two ways for retrieving a profile. You can use either `short_id` or `nickname`

**Request**

Query params available

1. **`sid`**: this query enables search by **short_id**
2. **`nn`**: this query enables search by **nickname**

**Response**

Status: 201

```json
{
	"shortId": "abdExsk",
	"nickname": "schecoperez",
	"linkedin": "https://linkedin.com/@/schecoperez",
	"email": "schp@redbullracing.com",
	"whatsapp": null,
	"medium": null,
	"twitterX": null,
	"website": null,
	"createdAt": "0001-01-01T00:00:00Z",
	"updatedAt": null
}
```

