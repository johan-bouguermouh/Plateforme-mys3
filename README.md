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

### Intyégré les variables d'environnements au projet

Vous trouverez un `.env.exemple` à la racine du projet. Créez un fichier `.env` au même emplacement est assuréez-vous de reprendre les même termes. Ces variable d'environnements servent à la fois au lancement du projet mais aussi à sa containerisation.

```
# Need Minio settings:
S3_ENDPOINT="your-S3-endpoint"
S3_PORT= 9000
S3_ACCESSKEY= "your-S3-accesskey"
S3_SECRETKEY= "your-S3-secretkey"
S3_BUCKET= "your-minio-bucket"
DB_BOLT_PATH=my.db
```

> Si toutefois les variables ne son pas déclarées vous pouvez faire tourner le projet en local avec la command d'execution suivante, après l'installation des modules necessaires :

```powershell
go run app.go
```

Dans de telles cironcstances le programme assignera automatiquement le port `9000` ainsi que le nom et le path du store **bbolt** à la racine de `api-interface` sous le nom de `mydb`

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
