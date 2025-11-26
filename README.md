# Power 4 - Jeu Connect Four

Un jeu Connect Four basé sur le web construit avec Go, proposant le rendu côté serveur et une variante bonus avec des fonctionnalités personnalisées.

## Prérequis

- **Version Go** : `1.25.4`
- **Dernière version** : `v1.0.0`

## Preview
https://secret.osadeo.com

## Rendu côté serveur (Server-Side Rendering)

Ce projet utilise la bibliothèque standard de Go pour le rendu côté serveur :

- **`html/template`** : Pour l'analyse et l'exécution de modèles HTML
- **`net/http`** : Pour les fonctionnalités du serveur HTTP et la gestion des requêtes

Le serveur rend les modèles HTML à chaque requête, injectant les données d'état du jeu dans les modèles avant de les envoyer au client. Les modèles sont situés dans :
- `base/templates/` - Modèles du jeu de base
- `bonus/templates/` - Modèles de la variante bonus

## Routes

### Routes du jeu de base

| Méthode | Chemin | Handler | Description |
|---------|--------|---------|-------------|
| `GET` | `/health` | Vérification de santé | Retourne le statut "OK" |
| `GET` | `/` | `handlers.HomeHandler` | Page principale du jeu |
| `POST` | `/move` | `handlers.MoveHandler` | Gère le coup du joueur |
| `POST` | `/new-game` | `handlers.NewGameHandler` | Démarrer une nouvelle partie |
| `POST` | `/reset-scores` | `handlers.ResetScoresHandler` | Réinitialiser les scores des joueurs |

### Routes de la variante bonus

| Méthode | Chemin | Handler | Description |
|---------|--------|---------|-------------|
| `GET` | `/bonus` | Redirection | Redirige vers `/bonus/setup` |
| `GET` | `/bonus/setup` | `bonusHandlers.SetupHandler` | Page de configuration du jeu (surnoms et taille du plateau) |
| `POST` | `/bonus/start-game` | `bonusHandlers.StartGameHandler` | Initialiser le jeu avec des paramètres personnalisés |
| `GET` | `/bonus/game` | `bonusHandlers.GameHandler` | Page du jeu bonus |
| `POST` | `/bonus/move` | `bonusHandlers.MakeMove` | Gère le coup du joueur (avec gravité inversée) |
| `POST` | `/bonus/new-game` | `bonusHandlers.NewGameHandler` | Démarrer une revanche avec les mêmes paramètres |
| `POST` | `/bonus/reset-scores` | `bonusHandlers.ResetScoresHandler` | Réinitialiser les scores des joueurs |

## Fonctionnalités bonus

La variante bonus inclut :
- **Surnoms des joueurs** : Noms personnalisés pour chaque joueur
- **Taille de plateau personnalisée** : Lignes et colonnes configurables (4-15)
- **Gravité inversée** : Tous les 5 coups, la gravité s'inverse (les pièces tombent du bas vers le haut)

## Lancement du serveur

```bash
go run main.go
```

Le serveur démarrera sur `127.0.0.1:80`

- Jeu de base : `http://127.0.0.1/`
- Variante bonus : `http://127.0.0.1/bonus/setup`

## Structure du projet

```
power4/
├── base/
│   ├── handlers/
│   │   └── handler.go      # Handlers du jeu de base
│   └── templates/
│       └── index.html      # Modèle du jeu de base
├── bonus/
│   ├── handlers/
│   │   └── handler.go      # Handlers de la variante bonus
│   └── templates/
│       ├── setup.html      # Modèle de la page de configuration
│       └── game.html       # Modèle du jeu bonus
├── shared/
│   ├── gamelogic.go        # Logique de jeu principale
│   └── server.go           # Configuration du serveur HTTP
├── main.go                 # Point d'entrée de l'application
└── go.mod                  # Définition du module Go
```

## Pile technologique

- **Backend** : Bibliothèque standard Go uniquement
  - `net/http` - Serveur HTTP et routage
  - `html/template` - Rendu de modèles côté serveur
- **Frontend** : HTML5 avec Tailwind CSS (via CDN)
