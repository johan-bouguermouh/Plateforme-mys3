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

## Utilisation des routes

### CreateBucket 

Route : 
http://127.0.0.1:3000/v1/bucket

Exemple de body (bucket complet) : 
<Bucket>
    <name>this-is-a-valid-bucket-name-1234</name>
    <creationDate>2024-09-11T12:34:56Z</creationDate>
    <owner>
        <ID>123456789</ID>
        <DisplayName>JohnDoe</DisplayName>
    </owner>
    <uri>https://example.com/mybucket</uri>
    <type>STANDARD</type>
    <storageClass>STANDARD</storageClass>
    <versioning>Enabled</versioning>
    <objectCount>12345</objectCount>
    <size>9876543210</size>
    <lastModified>2024-09-11T12:34:56Z</lastModified>
</Bucket>

Exemple de body CreateBucketRequest (partiel, suffisant pour create) : 
<Bucket>
    <name>this-is-a-valid-bucket-name-1234</name>
    <owner>
        <ID>123456789</ID>
        <DisplayName>JohnDoe</DisplayName>
    </owner>
    <type>STANDARD</type>
    <versioning>Enabled</versioning>
</Bucket>

Afin de créer un bucket, il faut veiller à respecter les règles de nommage AWS : 
Se référer à la partie ##1 du README middlewares

**UploadFiles**
Route :
POST /v1/bucket/:bucketName/upload

Description :
Télécharge un fichier dans le bucket spécifié.

Exemple de Body :
<UploadFileRequest>
  <FileName>example.txt</FileName>
  <FileContent>base64EncodedContent</FileContent>
</UploadFileRequest>

Réponse :
<UploadResponse>
  <message>Fichier téléchargé avec succès</message>
  <path>buckets/your-bucket-name/example.txt</path>
</UploadResponse>

**Notes**:

- FileContent doit être encodé en base64.
- Le middleware UploadFileValidationMiddleware valide le format XML et les champs requis avant que la requête ne soit traitée.

  **DeleteFile**
  Route :
 DELETE /v1/bucket/:bucketName/file/:fileName

Description :
Supprime un fichier du bucket spécifié.

Exemple de Response :
{
    "message": "Fichier supprimé avec succès"
}

**Notes**:

Le middleware FileExistenceMiddleware vérifie l'existence du fichier avant la suppression.

**ListFiles**
Route :
GET /v1/bucket/:bucketName/files

Description :
Liste tous les fichiers dans le bucket spécifié.

Exemple de Response :
<BucketObject>
    <Key>hero-bg.jpg</Key>
    <LastModified>2024-09-14T16:44:37+02:00</LastModified>
    <ETag></ETag>
    <Size>116112</Size>
    <StorageClass></StorageClass>
    <Owner>
        <UserKey></UserKey>
        <DisplayName></DisplayName>
        <Type></Type>
        <URI></URI>
        <ROLE>
            <Name></Name>
            <ID></ID>
            <Type></Type>
        </ROLE>
        <SecretKey></SecretKey>
    </Owner>
    <Type>OBJECT</Type>
    <URI>http://localhost:3000/buckets/johan/hero-bg.jpg</URI>
    <BucketName>johan</BucketName>
</BucketObject>

**Notes**:

Le middleware BucketExistenceMiddleware vérifie l'existence du bucket avant de lister les fichiers.
