FROM scratch
EXPOSE 3232
ADD search /
CMD ["/search"]
