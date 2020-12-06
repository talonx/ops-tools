# ops-tools
Various ops tools used in day-to-day operations on GCP and AWS

## iplocate
Usage: iplocate -project <gcp project> -ip <ip address> -region <region - unused right now>
  
  Attempts to determine if the IP address is owned by one of the running instances in the specific GCP project. Helpful for questions like "who is 123.34.56.78 who is making requests to my service?"
