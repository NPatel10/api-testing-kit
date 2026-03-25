FROM oven/bun:1.2.19 AS build

WORKDIR /workspace/apps/web

COPY apps/web/package.json apps/web/bun.lock ./
RUN bun install --frozen-lockfile

COPY apps/web/ ./

ARG PUBLIC_API_BASE_URL=http://localhost:8080
ENV PUBLIC_API_BASE_URL=$PUBLIC_API_BASE_URL

RUN bun run build

FROM oven/bun:1.2.19

WORKDIR /workspace/apps/web

COPY --from=build /workspace/apps/web /workspace/apps/web

ENV NODE_ENV=production

EXPOSE 4173

CMD ["bun", "run", "preview", "--", "--host", "0.0.0.0", "--port", "4173"]
