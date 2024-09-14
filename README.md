# Plateforme MYS3

## Initialiser le projet

### Prérequis

- Installer go
- Si utilisation VsCode installer le package Go _(Rich Go language support for Visual Studio Code)_
- Assurez-vous d'avoir Docer destop d'ouvert
- Au besoin de lancer le container api-interface tout seul, assurer vous d'avoir installer Make

  ```
  choco install make
  ```

  > _Assurez-vous de le lancer en mode administrateur_

### Cloner le repository

```
git remote add origin https://github.com/NouhaylaElfetehi/Plateforme-mys3.git
git branch -M main
git push -u origin main
```

### Installation des dépendances

```
# Clean packages
make clean-packages

# Generate go.mod & go.sum files
make requirements
```

Pour plus d'information rendez-vous sur le readme de `api-interface`

## Lancer le projet

### Lancer le projet en mode dev

Assurez vous que le package Air stable soit installer :

```
go install github.com/cosmtrek/air@v1.27.3
```

Lancer l'application en mode dev :

```
make start-dev
```

### Lancer le projet dans un container Docker

```
cd api-interface
make build
make up
marke start
```

### Lancer le projet en local

```
cd api-interface
go run app.go
```
