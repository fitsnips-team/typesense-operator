# Development Guide


## Docker

You will need to login to your container registry 

For Docker Hub run this command and enter username / password.

```
docker login registry-1.docker.io 
```

### Build and push docker image 
make docker-build docker-push  DOCKER_HUB_NAME=<container_registery_username>

```
typesense-operator on 🌱 main [📝🤷✓] via 🐹 v1.23.5 on 🅰 (us-west-2) took 11s 
🕙[ 02:15:25 ] ➜ make docker-build docker-push  DOCKER_HUB_NAME=jassinpain
docker build -t jassinpain/typesense-operator:0.2.12 .
[1/2] STEP 1/11: FROM golang:1.22 AS builder
[1/2] STEP 2/11: ARG TARGETOS
--> Using cache 1e0c3f51227b9ea4f3d306d4c7b8ef9bbeb0dfe6a5f08c5593a2b30b93176680
--> 1e0c3f51227b
[1/2] STEP 3/11: ARG TARGETARCH
--> Using cache 23c7fec145bebe78bbdd7a8fc81d896b7947c9c82c1a7cbc04a4ed3522d9280b
--> 23c7fec145be
[1/2] STEP 4/11: WORKDIR /workspace
--> Using cache 7bc83b61f5516fe0a3309f898e879f5ec9274b86588b83a807c132456ec00746
--> 7bc83b61f551
[1/2] STEP 5/11: COPY go.mod go.mod
--> Using cache 23aaadf02e70c1480b31abac7b82a1f9dc08588e052d3ff552a7841a23a47882
--> 23aaadf02e70
[1/2] STEP 6/11: COPY go.sum go.sum
--> Using cache 0a033d04f71fd79d0553776568a32eff4a15fe35494c845f449fae7ce2827d88
--> 0a033d04f71f
[1/2] STEP 7/11: RUN go mod download
--> Using cache dd0f2456f9a82d12f236452f7d4cb71e0652ce67c0d108ee6b17b8032337c53a
--> dd0f2456f9a8
[1/2] STEP 8/11: COPY cmd/main.go cmd/main.go
--> Using cache c81de85114dae48fbad1aeea291215c897b63925eb248ff32e588fc5e665bcd4
--> c81de85114da
[1/2] STEP 9/11: COPY api/ api/
--> Using cache f6a20a3ecd0547aa8f2d18ed9df8233852f7543dd77c4ead0ca44d3cd07396d4
--> f6a20a3ecd05
[1/2] STEP 10/11: COPY internal/controller/ internal/controller/
--> Using cache 0b904ffddaf39cf6f2f6105303404dbd03fd371d55acc543f1c01a5bccea2325
--> 0b904ffddaf3
[1/2] STEP 11/11: RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go
--> Using cache 21a9ad1a664b27dfda832b408e9f1ba093b9cd241379c8822724ee6839cb1d68
--> 21a9ad1a664b
[2/2] STEP 1/5: FROM gcr.io/distroless/static:nonroot
[2/2] STEP 2/5: WORKDIR /
--> Using cache 0555b91d53e63405539b9bbe0967de619dc7bab60f21326033526a2681e9ef70
--> 0555b91d53e6
[2/2] STEP 3/5: COPY --from=builder /workspace/manager .
--> Using cache 1249c42907e16f80e8d19e34da1c30f15a79ff1a153563a53c51b34f141fa2be
--> 1249c42907e1
[2/2] STEP 4/5: USER 65532:65532
--> Using cache 47576c3429f5552076afcd2d654e52a7a0f8ee36d5aec8fa806cd7396f537c5e
--> 47576c3429f5
[2/2] STEP 5/5: ENTRYPOINT ["/manager"]
--> Using cache 5bbcc4fdb7ef8a997c0a49265936a9ee09c73012b61ab7ec4c206333d0a3e825
[2/2] COMMIT jassinpain/typesense-operator:0.2.12
--> 5bbcc4fdb7ef
Successfully tagged localhost/jassinpain/typesense-operator:0.2.12
5bbcc4fdb7ef8a997c0a49265936a9ee09c73012b61ab7ec4c206333d0a3e825
docker push jassinpain/typesense-operator:0.2.12
Getting image source signatures
Copying blob sha256:6f1cdceb6a3146f0ccb986521156bef8a422cdbb0863396f7f751f575ba308f4
Copying blob sha256:8fa10c0194df9b7c054c90dbe482585f768a54428fc90a5b78a0066a123b1bba
Copying blob sha256:f920c5680b0b677741c5500dc365a9b074aa263ab43c3eb6be6a465b1caadd8e
Copying blob sha256:4d049f83d9cf21d1f5cc0e11deaf36df02790d0e60c1a3829538fb4b61685368
Copying blob sha256:a80545a98dcd0866ae5eeadc9a28dec703b1e54a01ce8ff245e83f48261fe575
Copying blob sha256:af5aa97ebe6ce1604747ec1e21af7136ded391bcabe4acef882e718a87c86bcc
Copying blob sha256:bbb6cacb8c82e4da4e8143e03351e939eab5e21ce0ef333c42e637af86c5217b
Copying blob sha256:2a92d6ac9e4fcc274d5168b217ca4458a9fec6f094ead68d99c77073f08caac1
Copying blob sha256:1a73b54f556b477f0a8b939d13c504a3b4f4db71f7a09c63afbc10acb3de5849
Copying blob sha256:f4aee9e53c42a22ed82451218c3ea03d1eea8d6ca8fbe8eb4e950304ba8a8bb3
Copying blob sha256:b336e209998fa5cf0eec3dabf93a21194198a35f4f75612d8da03693f8c30217
Copying blob sha256:41b3f0aff188e7619cde12306d17866bdbc0e67179c85923befa1c4a5931f2f4
Copying config sha256:5bbcc4fdb7ef8a997c0a49265936a9ee09c73012b61ab7ec4c206333d0a3e825
Writing manifest to image destination
```