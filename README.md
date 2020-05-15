# Vortodo

This application is a simple todo list application built as part of my interview with Vorto. The API uses the gRPC protocol, is written in Golang, uses PostgreSQL for a data layer, and is hosted on the Google Kubernetes Engine ([GKE](https://cloud.google.com/kubernetes-engine)), where it's IP address is exposed to public traffic. I have also written a restful client (also in golang) which can be used to interact with the gRPC back end; this client is hosted on the [App Engine](https://cloud.google.com/appengine). I initially used [Cloud Endpoints](https://cloud.google.com/endpoints) for the client, but my first efforts to deploy the full project were stymied by confusion around the authentication standards Endpoints demanded. I thus found it far simpler (both in terms of practical work and conceptually) to put the API on GKE with an exposed IP, then ping that IP via a restful client that can be touched via postman.

### Vortodo 1.0 background

Going into this project, I had no knowledge of gRPC, a passing familiarity with Kubernetes/GCloud, and a firm working knowledge of Go (though this knowledge was decidedly esoteric; I had little knowledge of how to build a web API like this). To be brief, I worked more or less nonstop from Friday, 5/1/2020 to Thursday, 5/7/2020. Despite initial successes, The efforts to stand up this Vortodo failed. The project timeline has a hinge at around Tuesday morning. Prior to that, I had spent quite a bit of time ascending the steep k8s/GCloud learning curve to get a beach head endpoint working. Unfortunately, at this point, I likely purged the credentials which allowed my client to communicate with the deployed server and lost the connection. This step backwards, coupled with a lack of sleep (and my refusal to sleep in the face of a looming deadline) caused me to spin my wheels and ultimately fail to produce a deliverable by the deadline.

### Vortodo 2.0

 After catching up on sleep, I opted to take a second crack at this project, reducing its complexity (primarily by removing the [micro](https://github.com/micro) library and by abandoning cloud endpoints) and extracting the core API code from the original to create the version I was actually able to stand up. The endpoints which exist at the time of writing (Creating a list, adding a task to a list, getting a list and its tasks, updating the completed status of a task and deleting a list with its tasks) were all more or less lifted from the wreckage of Vortodo 1.0.

 The design philosophy was basically gRPC pretending to be rest, in the interest of time. gRPC has such potential that I ultimately had to decide on a fairly restrictive series of endpoints in the name of forward progress; if I didn't, I knew I would be paralyzed by indecision in how to best approach this.

 #### Future Features

The most obvious future feature would be to implement user accounts and authentication, then refactor the API's core endpoints to work with the new user service.

A second line would be to continue to optimize the devops of the project. At present, the process of pushing new features live is fairly manual. I didn't really want to add the additional burden of learning CircleCI or Terraform; both would be reasonable next steps.

Third would be creating more endpoints for the list service itself and making it more flexible, really playing around with the potential of gRPC.

Finally, I considered the option of creating a demo which implemented two-way communication. My rough idea was that users could "trade lists with each other" or some other paper-thin excuse to play around with this capacity.  