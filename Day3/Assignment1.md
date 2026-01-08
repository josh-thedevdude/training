REQUEST TYPE        UNDER THE HOOD

GET                 The client made a get request which reaches the server where it goes through controllers,
                    services, repositories layers that does the business logic and data fetching from database.
                    Then it returns the response to the client with appropriate status code (200)

POST                The client makes a post request to create some data on the server by sending some body with
                    the request which goes through the same cycle on the server and creates new entry in the 
                    database and returns the created data with appropraite status code (201)

POSTMAN COLLECTION JSON:
{
	"info": {
		"_postman_id": "9cf1c3f8-1883-43eb-820b-2c4133358da0",
		"name": "BE_FUNDAMENTALS_A1",
		"description": "### Welcome to Postman! This is your first collection. \n\nCollections are your starting point for building and testing APIs. You can use this one to:\n\n• Group related requests\n• Test your API in real-world scenarios\n• Document and share your requests\n\nUpdate the name and overview whenever you’re ready to make it yours.\n\n[Learn more about Postman Collections.](https://learning.postman.com/docs/collections/collections-overview/)",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		"_exporter_id": "51337727",
		"_collection_link": "https://sahil-patil-34033428-4218891.postman.co/workspace/Sahil-Patil's-Workspace~aca121b4-5a6d-4a0b-86d4-21ceb603eb2b/collection/51337727-9cf1c3f8-1883-43eb-820b-2c4133358da0?action=share&source=collection_link&creator=51337727"
	},
	"item": [
		{
			"name": "Pokemon Search",
			"event": [
				{
					"listen": "test",
					"script": {
						"exec": [
							"pm.test(\"Status code is 200\", function () {",
							"    pm.response.to.have.status(200);",
							"});"
						],
						"type": "text/javascript",
						"packages": {},
						"requests": {}
					}
				}
			],
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "https://pokeapi.co/api/v2/pokemon/pikachu",
					"protocol": "https",
					"host": [
						"pokeapi",
						"co"
					],
					"path": [
						"api",
						"v2",
						"pokemon",
						"pikachu"
					]
				},
				"description": "This is a GET request and it is used to \"get\" data from an endpoint. There is no request body for a GET request, but you can use query parameters to help specify the resource you want data on (e.g., in this request, we have `id=1`).\n\nA successful GET response will have a `200 OK` status, and should include some kind of response body - for example, HTML web content or JSON data."
			},
			"response": []
		},
		{
			"name": "create JSON data",
			"event": [
				{
					"listen": "test",
					"script": {
						"exec": [
							"pm.test(\"Successful POST request\", function () {",
							"    pm.expect(pm.response.code).to.be.oneOf([200, 201]);",
							"});",
							""
						],
						"type": "text/javascript",
						"packages": {},
						"requests": {}
					}
				}
			],
			"request": {
				"method": "POST",
				"header": [],
				"body": {
					"mode": "raw",
					"raw": "{\n  \"name\": \"sahil patil\",\n  \"assignment\": 1,\n  \"description\": \"Backend Fundamentals\"\n}",
					"options": {
						"raw": {
							"language": "json"
						}
					}
				},
				"url": {
					"raw": "https://jsonplaceholder.typicode.com/posts",
					"protocol": "https",
					"host": [
						"jsonplaceholder",
						"typicode",
						"com"
					],
					"path": [
						"posts"
					]
				},
				"description": "This is a POST request, submitting data to an API via the request body. This request submits JSON data, and the data is reflected in the response.\n\nA successful POST request typically returns a `200 OK` or `201 Created` response code."
			},
			"response": []
		}
	]
}