# Usage d'un Modèle Repository

## Avant Propos

### Usage et Obligation

Le programme ne necéssite pas forcement l'intégration d'un Modèle Repositoy pour fonctionner correctement. Ce dernier offre seulement un couche d'abstraction, comme un boite à outil pour faciliter le travail des devs. Libre à chacun de s'en apparait. Cependant, son usage permet de s'assurer de l'intégrité des données et de respecté quelques normes.

### Utilité

La documentation suivante permet de de définire, instancier et exploité une nouvelle entitée liées dynamiquement a bbold. Cette entoté pourra pas la suite exploité une QueryBuilder permattant d'ajouté un couche d'abstraction à la gestion de donnée dans `bbolt` .

### Intégration de Bbolt comme store

[Bbolt](https://github.com/etcd-io/bbolt/tree/main?tab=readme-ov-file#getting-started) est un fork du magasin de clés/valeurs Bolt [de Ben Johnson](https://github.com/benbjohnson) . Ce store de clés/valeurs purement Go fournit une base de données simple, rapide et fiable pour les projets qui ne nécessitent pas de serveur de base de données complet tel que Postgres ou MySQL. A travers son usage, nous chercherons à optimiser la recherche d'élément au sain de notre programme.

## Prérequis

## S"assurer de la présence du paquet Bbolt

Pour utiliser cet ORM, vous devez d'abord installer BoltDB. Vous pouvez l'ajouter à votre projet Go en utilisant `go get` :

```bash
go get go.etcd.io/bbolt
```

### Initialisation de la Base de Données

Avant d'utiliser l'ORM, assurez vous que la base de données BoltDB et les repositories sont présents :

```go
import (
    "api-interface/models"
    "go.etcd.io/bbolt"
)

func main() {
    db, err := bbolt.Open("mydb.db", 0600, nil) //initialise la base de donnée Bolt
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    models.InitRepositories(db) //Initialise les Repository Selon les entités
}
```

> ATTENTION : la méthode `InitRepositories` ne doit être appeller que si des entités sont déclarés

### Déclaration des Entités

Les entités doivent implémenter l'interface [`Entity`](https://github.com/NouhaylaElfetehi/Plateforme-mys3/blob/339ac753efa7beb8acb0e398834f318c1fe3a8a2/api-interface/entities/entity.interface.go):

#### Déclarer une nouvelle entité dans l'interface entity

Nous prendrons ici comme exemple l'Entity `Bucket` comme exemple

```go
package entity

type Bucket struct {
    Name string `json:"name"`
    // autres champs...
}

func (b *Bucket) GetKey() string {
    return b.Name
}

func (b *Bucket) Serialize() ([]byte, error) {
    return json.Marshal(b)
}

func (b *Bucket) Deserialize(data []byte) error {
    return json.Unmarshal(data, b)
}
```

Pour l'instant l'interface comprend seulement trois méthodes :

- **GetKey** : Retourne l'attribut aillant été designer comme clef primaire de notre type d'object qui sera servi dans notre store
- Serialize : Fonction qui permet de Serialiser les donnée en bytes
- Deserialize : Qui permet de les extraires et les interprétés

#### Conventions attendus

J'ai choisi comme de suivre la **convention de nommage à base de suffixes** ou **convention de nommage modulaire**. Cette approche nous permet de définir les type de fihcier a lire selon leurs particularité. Cela nous permet de parcourir les entités indépendament des autres fichiers pouvant être présent dans un package donné.

Je conseille à toutes et tous de poursuivre cette approche pour un meilleur intégration du code concernant se module.

### Déclaration et initialisation du Repository

### AutoInstacier votre Repository

Enregistrez vos entités dans le fichier [`models.go`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html) :

```go
func init() {
    Use[*entity.Bucket]()
    // autres entités...
}
```

Initialise les repositories pour chaque entité déclarée utilisé dans la fonction init() de ce fichier. Cette fonction est appelée au démarrage de l'application pour initialiser les repositories pour chaque entité déclarée. Init repositories utilise les QueryBuilder pour chaque entité déclarée et les stocke dans un map pour une utilisation ultérieure.

#### Contexte de [`InitRepositories`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition") et [`NewQueryBuilder`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition")

Lors de l'initialisation des repositories avec [`InitRepositories`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition"), vous créez déjà des instances de [`QueryBuilder`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition") pour chaque type d'entité et les stockez dans une map globale [`repositories`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition"). Chaque [`QueryBuilder`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition") est initialisé avec la base de données ([`db`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition")) qui créer directement un bucket bold pour prêt à être consommé. Le bucket sera nommé au pluriel et selon le nom de votre entitée

#### Utilisation de [`repositories["Bucket"]`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition")

Lorsque vous récupérez [`repositories["Bucket"]`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition"), vous obtenez une instance de [`QueryBuilder`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html "Go to definition") qui est déjà configurée avec la base de données et prête à l'emploi. Cela signifie que vous n'avez pas besoin de réinstancier ou de reconfigurer la base de données dans votre modèle.

## Utilisation

### Créer un Modèle implementant le repositoy

```go
package models

import (
    entity "api-interface/entities"
    "api-interface/repositories"
)

type BucketModel struct {
    bucketRepository *repositories.QueryBuilder[*entity.Bucket]
}

//Le queryBuilder Peut être déclaré dans la function init()
func UseBucketModel() (*BucketModel, error) {
    queryBuilder, err := UseRepository[*entity.Bucket]("Bucket")
    if err != nil {
        return nil, err
    }

    return &BucketModel{
        bucketRepository: queryBuilder,
    }, nil
}

func (bm *BucketModel) Insert(bucket *entities.Bucket) error {
    return bm.bucketRepository.Insert(bucket)
}
```

Utilisez les méthodes du [`QueryBuilder`](vscode-file://vscode-app/c:/Users/laplateforme/AppData/Local/Programs/Microsoft%20VS%20Code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html) pour effectuer des opérations CRUD :

```go
func main() {
    bucketModel, err := models.UseBucketModel()
    if err != nil {
        log.Fatal(err)
    }

    newBucket := &entities.Bucket{Name: "example-bucket"}
    err = bucketModel.Insert(newBucket)
    if err != nil {
        log.Fatal(err)
    }

    // Récupérer un bucket
    var bucket entities.Bucket
    err = bucketModel.bucketRepository.Get("example-bucket", &bucket)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Bucket récupéré :", bucket)
```

## Methodes QueryBuilder

| Méthode                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | Type de Retour                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [`GetEntity`](command:_github.copilot.openSymbolFromReferences?%5B%22%22%2C%5B%7B%22uri%22%3A%7B%22scheme%22%3A%22file%22%2C%22authority%22%3A%22%22%2C%22path%22%3A%22%2Fc%3A%2FUsers%2Flaplateforme%2FDocuments%2FGitHub%2FPlateforme-mys3%2Fapi-interface%2Frepositories%2Fquery_builder.go%22%2C%22query%22%3A%22%22%2C%22fragment%22%3A%22%22%7D%2C%22pos%22%3A%7B%22line%22%3A19%2C%22character%22%3A27%7D%7D%5D%2C%22f44afa0f-a116-42dc-bdc9-702f7187b5cc%22%5D "Go to definition") | [`T`](command:_github.copilot.openSymbolFromReferences?%5B%22%22%2C%5B%7B%22uri%22%3A%7B%22scheme%22%3A%22file%22%2C%22authority%22%3A%22%22%2C%22path%22%3A%22%2Fc%3A%2FUsers%2Flaplateforme%2FDocuments%2FGitHub%2FPlateforme-mys3%2Fapi-interface%2Frepositories%2Fquery_builder.go%22%2C%22query%22%3A%22%22%2C%22fragment%22%3A%22%22%7D%2C%22pos%22%3A%7B%22line%22%3A10%2C%22character%22%3A18%7D%7D%5D%2C%22f44afa0f-a116-42dc-bdc9-702f7187b5cc%22%5D "Go to definition") | Retourne l'entité associée au[`QueryBuilder`](command:_github.copilot.openSymbolFromReferences?%5B%22%22%2C%5B%7B%22uri%22%3A%7B%22scheme%22%3A%22file%22%2C%22authority%22%3A%22%22%2C%22path%22%3A%22%2Fc%3A%2FUsers%2Flaplateforme%2FDocuments%2FGitHub%2FPlateforme-mys3%2Fapi-interface%2Frepositories%2Fquery_builder.go%22%2C%22query%22%3A%22%22%2C%22fragment%22%3A%22%22%7D%2C%22pos%22%3A%7B%22line%22%3A9%2C%22character%22%3A3%7D%7D%5D%2C%22f44afa0f-a116-42dc-bdc9-702f7187b5cc%22%5D "Go to definition").                        |
| `NewQueryBuilder`                                                                                                                                                                                                                                                                                                                                                                                                                                                                          | `(*QueryBuilder[T], error)`                                                                                                                                                                                                                                                                                                                                                                                                                                                        | Crée une nouvelle instance de[`QueryBuilder`](command:_github.copilot.openSymbolFromReferences?%5B%22%22%2C%5B%7B%22uri%22%3A%7B%22scheme%22%3A%22file%22%2C%22authority%22%3A%22%22%2C%22path%22%3A%22%2Fc%3A%2FUsers%2Flaplateforme%2FDocuments%2FGitHub%2FPlateforme-mys3%2Fapi-interface%2Frepositories%2Fquery_builder.go%22%2C%22query%22%3A%22%22%2C%22fragment%22%3A%22%22%7D%2C%22pos%22%3A%7B%22line%22%3A9%2C%22character%22%3A3%7D%7D%5D%2C%22f44afa0f-a116-42dc-bdc9-702f7187b5cc%22%5D "Go to definition") pour une entité donnée. |
| `getBucket`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                | `(*bbolt.Bucket, error)`                                                                                                                                                                                                                                                                                                                                                                                                                                                           | Retourne le bucket associé au[`QueryBuilder`](command:_github.copilot.openSymbolFromReferences?%5B%22%22%2C%5B%7B%22uri%22%3A%7B%22scheme%22%3A%22file%22%2C%22authority%22%3A%22%22%2C%22path%22%3A%22%2Fc%3A%2FUsers%2Flaplateforme%2FDocuments%2FGitHub%2FPlateforme-mys3%2Fapi-interface%2Frepositories%2Fquery_builder.go%22%2C%22query%22%3A%22%22%2C%22fragment%22%3A%22%22%7D%2C%22pos%22%3A%7B%22line%22%3A9%2C%22character%22%3A3%7D%7D%5D%2C%22f44afa0f-a116-42dc-bdc9-702f7187b5cc%22%5D "Go to definition"). (méthode privée)       |
| `Insert`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | `error`                                                                                                                                                                                                                                                                                                                                                                                                                                                                            | Insère une entité dans la base de données.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| `Get`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      | `error`                                                                                                                                                                                                                                                                                                                                                                                                                                                                            | Récupère une entité de la base de données en utilisant sa clé.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| `GetAll`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | `([]T, error)`                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | Récupère toutes les entités de la base de données qui passent le filtre donné.                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |

## Contribuer

Les contributions sont les bienvenues ! Veuillez soumettre des pull requests ou ouvrir des issues pour signaler des bugs ou proposer des améliorations. Le QueryBuilder est encore à l'étape préliminaire. Les méthodes doivent être enrichies.

> Ce README devrait fournir une base solide pour les développeurs souhaitant utiliser et étendre votre ORM customisé. Assurez-vous de le mettre à jour au fur et à mesure que de nouvelles fonctionnalités sont ajoutées.
> Créer par Johan Bouguermouh
