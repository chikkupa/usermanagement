# Standard Google Cloud Platform
enter directory code and run `dev_appserver.py app.yaml`

# Docker version
## win
Install docker from https://www.docker.com/
## mac
Install docker from https://www.docker.com/
## linux
install docker and docker-compose from package manager. Latest versions not available via package managers.
Install directly from site for latest version

## Check installation
docker installation. I am running version 18.05.0-ce
`docker -v`

docker compose installation. minimum version required 1.6.0+. I am running version 1.21.2
note: downgraded docker-compose file to version '2' from '3.1' so as to support docker-compose version 1.6.0+
`docker-compose -v`

## Configure Google Cloud
Pull latest google/cloud-sdk image
`docker pull google/cloud-sdk:latest`

Configure Cloud
`docker run -ti --name gcloud-config google/cloud-sdk gcloud auth login`

## Run container
`docker-compose up`

Docker shall install required images (software), create and launch a container

The container is configured to delete all data stored within it on start thus providing an easy dev setup


### Test users

In order to test account status control 3 pages have been created /defaultPage , /premiumPage and /adminPage

Please use the `Is User <user type>?` link in the menu to test this functionality

Default User: test@example.com
Can access: /defaultPage
No access: /premiumPage
No access: /adminPage

Premium User: premium@example.com
Can access: /premiumPage
No access: /adminPage
No access: /defaultPage

Admin User: admin@example.com
Can access: /adminPage
No access: /defaultPage
No access: /premiumPage

These users are created on the first login by a user, and are deleted once the applicaton is shutdown.
More details of the users can be seen on their profile page.

### Test Virtual Applications

In order to add test the users we have created 3 test endpoints to be used to emulate applications.
/testapp01
/testapp02
/testapp03

These endpoints correspond to the names of the expected test applications.
At the moment only one, testapp01, is created by default on login.
The other applications have not been created so that the create page functionality can be tested.
To test the create app functionality please navigate to the create new application nav option in the menu.
Just change the ApplicationID, which is the primary key, to testapp02 or testapp03 and submit to create the appliations.

Once the test application has been created you can navigate to the page and add the application to a user. By navigating to say /testapp01 and submiting a user mail address the users applications column is populated with a new application id. This id cannot be entered twice and is unique.

Looking forward to hearing on how to move forward with this. Please let me know if the above is as expected.

## Details

Usermanagement service is running on http://localhost:8080.
Once started the service shall not contain any data. This allows us to easily destroy and restart the service with a fresh database

/* Deprecated*/
POCtest service emulates inserting data after user signin.
Simply navigate to http://localhost:9090/ and signin to the service.
To verify the insertion of data visit http://localhost:8000/datastore
