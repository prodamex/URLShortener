# Simple Go URL Shortener

Ce projet est un système de raccourcissement d'URL simple mais fonctionnel écrit en Go, utilisant MongoDB pour stocker les correspondances entre les URL longues et courtes. Il offre une redirection efficace de l'URL courte vers l'URL longue et fournit des statistiques basiques, telles que le nombre de liens raccourcis et le nombre de clics par lien.

## Fonctionnalités

- **Raccourcissement d'URL**: Génère des URL courtes uniques pour des URL longues.
- **Redirection**: Redirige les utilisateurs vers l'URL longue lorsqu'ils visitent l'URL courte.
- **Statistiques**: Affiche le nombre de liens raccourcis et le nombre de clics sur chaque lien.
- **Base de données MongoDB**: Utilise MongoDB pour stocker les relations entre les URL courtes et longues.

## Prérequis

Pour exécuter ce projet, vous aurez besoin de :

- Go (version 1.x)
- MongoDB installé et en cours d'exécution sur votre machine

## Installation

Clonez le dépôt sur votre machine locale en utilisant :

```bash
git clone 
```

Naviguez dans le répertoire du projet :
```
cd ..
```

##Configuration
Avant de lancer l'application, assurez-vous que MongoDB est installé et en cours d'exécution. Vous pouvez configurer la chaîne de connexion MongoDB dans un fichier de configuration ou directement dans votre code.

##Exécution
Pour démarrer le serveur sur le port 3030, exécutez :
```
go run main.go
```
**L'application devrait maintenant être en cours d'exécution et écouter sur le port 3030. Vous pouvez accéder à l'API de raccourcissement d'URL via http://localhost:3030.**

## Utilisation
Pour raccourcir une URL, envoyez une requête HTTP POST à /shorten avec l'URL longue dans le corps de la requête. Pour accéder à une URL raccourcie, naviguez simplement vers l'URL courte générée, et vous serez redirigé vers l'URL longue originale.

