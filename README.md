# About
Fooder was created for domestic usage to ease meal planning and grocery shopping.  
This is a REST api created in Golang - as vanilla as possible - and serves me as personal playground and proof of Golang knowledge.

# Installation
To run this project either compile for your machine architecture or build Docker image. GOARCH build env might need to be adjusted to your requirements at the time of writting this README, you'd need to edit the Dockerfile.

# Architecture
I'm a big fan of modular monoliths as a base even tho this project started as one file and had been refactored into modules. 
Architecture of this project is evolutionary one. 
I'm not using any frameworks like Gorilla or Gin
but my own "framework" packages can be found in `/internal/framework`  
"Domain" packages follow simple layered structure. There are no explicit `Application|Domain|Infrastructure` layers,
since only I work on this project. Layers are implied.  
* service - Application
* model - Domain
* x_repository - Infrastructure
and so on.
For the simplicity I started with SQLite DB but I'm thinking on transitioning to Postgres for its robustness or Mongo(for educational purposes)