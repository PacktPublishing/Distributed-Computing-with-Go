default: clean start

clean:
	@echo "Delete deployments..."
	-kubectl delete deployments --all
	@echo "Done."
	@echo "Deleting services..."
	-kubectl delete services file-server n-z-librarian others-librarian a-m-librarian concierge
	@echo "Done."
	@echo

start:
	@echo "Ensuring goophr-config exists..."
	-kubectl create cm goophr-config --from-file=./config-map.yaml
	@echo "Creating services & deployments..."
	kubectl create -f concierge.yaml -f others_librarian.yaml -f a_m_librarian.yaml -f n_z_librarian.yaml -f file-server.yaml

