FROM teamfruit/nginx-fcgiwrap:latest

# Run updates
RUN apt update && apt install git fcgiwrap apache2-utils nginx -y

RUN mkdir -p /srv/git && mkdir -p /srv/git-evil

# Setup git server
RUN git config --global --bool core.bare true && git config --global user.email "you@example.com" && git config --global user.name "Your Name"
RUN git init /srv/git/foo --bare --shared && git -C /srv/git/foo update-server-info && git -C /srv/git/foo config --local http.receivepack true
RUN git init /srv/git-evil/foo --bare --shared && git -C /srv/git-evil/foo update-server-info && git -C /srv/git-evil/foo config --local http.receivepack true

# Generate git repo
RUN git clone /srv/git/foo /tmp/foo && git clone /srv/git-evil/foo /tmp/foo-evil
COPY go.mod /tmp/foo/
COPY go.mod /tmp/foo-evil/
COPY normal.go /tmp/foo/
COPY evil.go /tmp/foo-evil/
RUN cd /tmp/foo && git add . && git commit -m init && git push origin master && git tag -m 'v1.0.0' v1.0.0 && git push --tags
RUN cd /tmp/foo-evil && git add . && git commit -m init && git push origin master && git tag -m 'v1.0.0' v1.0.0 && git push --tags
RUN rm -rf /tmp/foo /tmp/foo-evil

# Give that directory the proper ownership for nginx
RUN chown -R www-data:www-data /srv/git && chmod -R 755 /srv/git
RUN chown -R www-data:www-data /srv/git-evil && chmod -R 755 /srv/git-evil

COPY default.conf /etc/nginx/conf.d/.
