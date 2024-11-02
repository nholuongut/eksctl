FROM scratch
LABEL maintainer="Nho Luong <luongutnho@hotmail.com>"
CMD eksctl
COPY --from=nholuongut/eksctl-builder:latest /out /
