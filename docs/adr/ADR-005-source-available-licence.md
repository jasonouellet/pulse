# ADR 005 : Modèle de Licence Source-Available avec Restriction Commerciale

## Statut

Accepté

## Contexte

Le Project PULSE vise à fournir une solution moderne de gestion de sports collectifs. Nous souhaitons encourager l'adoption communautaire, la transparence du code et permettre aux clubs sportifs à but non lucratif de s'auto-héberger gratuitement.

Cependant, les licences open-source permissives (MIT, Apache 2.0) ou Copyleft standards (GPLv3, AGPLv3) ne protègent pas suffisamment l'investissement commercial de l'auteur principal. Elles permettent à des fournisseurs de services Cloud tiers ou des agences de prendre le code et de le commercialiser en SaaS sans compensation.

## Décisions

1. **Licence Source-Available Non-Commercial (BSL 1.1 / PolyForm NonCommercial) :**
   * **Auto-hébergement par les clubs :** Tout club sportif ou association à but non lucratif a le droit de télécharger et déployer l'application gratuitement pour sa propre gestion interne.
   * **Transparence :** Le code source reste public, auditable et ouvert aux contributions.
2. **Restrictions Commerciales :**
   * Il est strictement interdit à un tiers de proposer le Project PULSE sous forme de service hébergé payant (SaaS / Cloud Managed) ou de revendre le logiciel.
3. **Modèle de Double Licence (Dual-Licensing) :**
   * L'auteur principal conserve 100 % du Copyright. Une entreprise ou une ligue professionnelle souhaitant exploiter le logiciel dans un cadre commercial dérogeant à la licence doit acquérir une **Licence Commerciale** auprès de l'auteur.

## Conséquences

* Protection absolue contre la captation de valeur commerciale par des tiers.
* Auto-hébergement totalement gratuit garanti pour les clubs amateurs.
* Exigence d'un accord de contribution (CLA.md) pour les contributeurs externes.
