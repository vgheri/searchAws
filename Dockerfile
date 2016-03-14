FROM scratch
EXPOSE 3232
ADD searchAws /
CMD ["/searchAws"]
