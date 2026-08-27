# ADR-006 — Obligations de télémétrie et de santé

## Décision

Tout service déployable doit exposer une sonde de liveness indépendante de ses dépendances, une sonde de readiness qui vérifie ses dépendances critiques et émettre sa télémétrie via OpenTelemetry/OTLP. Les configurations d'endpoint doivent être injectées par l'environnement, jamais codées avec des secrets.

## Conséquences

Les Deployments Kubernetes définissent des probes `startup`, `liveness` et `readiness`, ainsi que des limites de ressources. Tout nouveau service doit ajouter la documentation de ses endpoints de santé et être compatible avec un Collector OTLP. Les erreurs de santé ne doivent pas contenir de secret.
