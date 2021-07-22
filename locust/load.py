from locust import HttpUser, between, task


class WebsiteUser(HttpUser):
    wait_time = between(5, 15)



    @task
    def keyword(self):
        url = '/createGift'
        data = {"codeType":"3","drawCount":"10","des":"this des","validTime":"1","content":"{1:1000}","createUser":"qq","userId":""}
        self.client.post(url=url,data=data)
        self.client.get("/getGift?code=jG7a4lo8")
        url = '/checkCode'
        data = {"90KKHauh":"888888"}
        self.client.post(url=url,data=data)