export KO_DOCKER_REPO=eu.gcr.io/herman-codam-tmp-nov-2022
calculator=$(ko build ./calculator) && \
equations=$(ko build ./equations) && \
frontend=$(ko build ./)

gcloud run deploy calculator -q --image $calculator --region=europe-west4 --allow-unauthenticated --set-env-vars=OTEL_SERVICE_NAME=calculator,OTEL_RESOURCE_ATTRIBUTES=g.co/gae/app/module=calculator
gcloud run deploy equations -q --image $equations --region=europe-west4 --allow-unauthenticated --set-env-vars=OTEL_SERVICE_NAME=equations,OTEL_RESOURCE_ATTRIBUTES=g.co/gae/app/module=equations

url_calculator=$(gcloud run services describe calculator --platform managed --region=europe-west4 --format 'value(status.url)')
url_equations=$(gcloud run services describe equations --platform managed --region=europe-west4 --format 'value(status.url)')

gcloud run deploy frontend -q --image $frontend --region=europe-west4 --allow-unauthenticated --set-env-vars=URL_EQUATIONS=$url_equations,URL_CALCULATOR=$url_calculator,OTEL_SERVICE_NAME=frontend,OTEL_RESOURCE_ATTRIBUTES=g.co/gae/app/module=frontend
