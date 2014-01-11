./build.sh

if [ -e "goracle.id" ] 
then
	echo Removing existing container $ID
	ID=(cat goracle.id)
	docker kill $ID
	docker rm $ID
fi

docker run -d goracle-standalone > goracle.id
