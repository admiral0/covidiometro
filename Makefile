covidiometro_arm:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o covidiometro_arm

covidiometro.it: covidiometro_arm
	scp covidiometro_arm pi@covidiometro.it:/home/pi/covidiometro_new

clean:
	rm covidiometro_arm

.PHONY: covidiometro_arm covidiometro.it