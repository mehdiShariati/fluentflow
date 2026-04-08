# Deployment

This page explains how to **host the repository on GitHub**, publish **online documentation** with **GitHub Pages**, and outline **production deployment** patterns for the FluentFlow stack.

## 1. Publish the project on GitHub

1. Create a new repository on GitHub (empty, no README if you are pushing an existing tree).
2. On your machine, from the project root:

   ```bash
   git init
   git add .
   git commit -m "Initial import: FluentFlow stack"
   git branch -M main
   git remote add origin https://github.com/YOUR_ORG/YOUR_REPO.git
   git push -u origin main
   ```

3. Replace `YOUR_ORG/YOUR_REPO` with your account or organization and repository name.

4. **Protect `main`** (optional but recommended): Settings → Branches → add a rule requiring pull requests and status checks before merge.

## 2. Online documentation with GitHub Pages

This repository includes **MkDocs** configuration ([`mkdocs.yml`](https://github.com/mehdi/fluentflow/blob/main/mkdocs.yml)) and a workflow **`.github/workflows/docs.yml`** that builds the site and deploys it to **GitHub Pages**.

### Enable Pages

1. In the GitHub repository: **Settings** → **Pages**.
2. Under **Build and deployment**, set **Source** to **GitHub Actions** (not “Deploy from a branch” unless you prefer a legacy flow).

### Deploy

- Every **push** to **`main`** or **`master`** runs **Deploy documentation** (MkDocs → GitHub Pages). You do **not** need to touch only `docs/` for the first deploy.
- You can also run it by hand: **Actions** → **Deploy documentation** → **Run workflow** (pick your default branch).

### If `https://<user>.github.io/<repo>/` returns 404

Usually nothing has been published yet, or the last workflow **failed**.

1. Open **Actions** → **Deploy documentation** and confirm a **green** run on your default branch.
2. If there is no run, push this repo (including `.github/workflows/docs.yml`) or trigger **Run workflow** manually.
3. Use the **MkDocs** workflow only. Ignore GitHub’s suggested **Next.js** Pages template—this site is **MkDocs**, not a Next.js static export.
4. **Public repos** get free GitHub Pages; **private** repos need a paid GitHub plan for Pages in many cases.

### Where to find the live docs

After the first successful run, the site URL is:

```text
https://YOUR_GITHUB_USERNAME.github.io/YOUR_REPO_NAME/
```

If the repository is under an organization, use that hostname segment instead of your username.

### Optional: custom domain

1. Add a `CNAME` file to the `docs/` folder (MkDocs copies it to the site root) containing your domain name, or configure the domain in **Settings → Pages**.
2. Configure DNS `A`/`CNAME` records as described in [GitHub Pages custom domains](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site).

### Local preview of documentation

```bash
pip install -r requirements-docs.txt
mkdocs serve
```

Open **http://127.0.0.1:8000** to preview before pushing.

## 3. Deploying the application (production overview)

FluentFlow is a **multi-service** system: **PostgreSQL**, **LiveKit**, **Go API**, **Next.js**, and a **Python agent worker**. Production deployment is not a single static binary; you typically orchestrate these with a container platform or VMs.

### Environment variables (high level)

| Area | Examples |
|------|----------|
| API | `DATABASE_URL`, `JWT_SECRET`, `CORS_ORIGINS`, `LIVEKIT_*`, `OPENAI_API_KEY`, optional `ADMIN_TOKEN` |
| Web | `NEXT_PUBLIC_API_URL` must point to the public API URL learners use |
| Agent | `LIVEKIT_*`, `OPENAI_*`, `LIVEKIT_AGENT_NAME` matching API-issued tokens |

Never commit real secrets; use your host’s secret store (GitHub Actions secrets, AWS Secrets Manager, etc.).

### Suggested hosting patterns

- **Containers:** Build images from `Dockerfile.api`, `web/Dockerfile`, and `agent/Dockerfile`; run with Kubernetes, ECS, Nomad, or Docker Swarm.
- **Managed LiveKit:** [LiveKit Cloud](https://livekit.io/cloud) reduces SFU operations; point `LIVEKIT_URL` and API credentials at the cloud project.
- **Database:** Managed PostgreSQL (RDS, Cloud SQL, Azure Database, etc.) with TLS and automated backups.
- **CI/CD:** Extend `.github/workflows` with jobs that build and push images to a registry and deploy to your environment (not included by default — environment-specific).

### HTTPS and the browser

The learner app uses **WebRTC**; production must serve the **web app over HTTPS** and use **WSS** for LiveKit. Terminate TLS at a load balancer or reverse proxy (Caddy, nginx, cloud LB).

### GitHub Actions and “deploy”

**GitHub Pages** in this repo deploys **documentation only**. Application deployment to your cloud is usually a separate pipeline: build images, run migrations, roll out services. The [Scaling](scaling.md) page discusses operational concerns that inform that pipeline.
