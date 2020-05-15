# gcloud container clusters create simpletodo \
#    --num-nodes 2 \
#    --enable-basic-auth \
#    --issue-client-certificate \
#    --region us-central1

gcloud builds submit --tag gcr.io/vortodo/simpletodo .