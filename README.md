Ubuntu
------

1. Clone Patrao Upgrade service repository (https://github.com/nutanix/patrao)

2. Install golang package  1.12.1 or greater 
         - sudo snap install --classic go

3. Open <patrao_root>/cmd/upgradeagent/ folder in Terminal and run the next commands:
    - go get
    - go build   <- just for make sure that all packages has been downloaded

4. Install docker service (sudo apt install docker.io)
    - after installation add your user to docker group (sudo usermod -aG docker $USER)
    - reboot your host system 

5. Install Visual Studio Code (download link:  https://go.microsoft.com/fwlink/?LinkID=760868)


6. Open Visual Studio Code and go to menu  âFile/Open Workspaceâ (<patrao_root>/cmd/Patrao.code-workspace)

7. Download and install mock-server (www.mock-server.com) 
    - start mock server like (mockserver -serverPort 1080

8. Run mock setup script <patrao_root>/scripts/upgrade_enigne/setup_mock.sh (might be need to install curl on your system)

9. Create docker-compose.yml with next: 

--------
version: "2"
services:
  db:
    image: postgres:10.3
    expose:
      - 5432
    environment:
      - POSTGRES_USER=boosteroid
      - POSTGRES_PASSWORD=P56FJXce6P60QQMa7QqNbEf4Z1Mm6uBE
  cache: 
    image: postgres:9.5
----

Copy docker-compose.yml under âtestâ folder and launch test solution (docker-compose up -d)

10. Press Ctrl+Shif+B and run the next sequence of commands: 
   - Open <patrao_root>/deployments/patraoagent/Dockerfile and change this line
 
CMD ["./upgradeagent", "--upstreamHost=http://192.168.1.58:1080"]

Like this: 
CMD ["./upgradeagent", "--upstreamHost=http://localhost:1080"]


   - 1. upgrade_agent_build_and_deploy.   (Build solution)
   - 2. upgrade_agent_create_image (create docker image with Patrao agent on board)
   - 3. upgrade_agent_run_image  (runs created image as docker container)

After that Patrao Agent should be launched under docker container and in 30 sec will do upgrade attempt. Mock server return to agent a new specification so launched 
"test" solution will be upgraded. On the next round of upgrade (in 30 sec) Patrao agent detect that âtestâ solution are up to date and wonât do upgrade. 