FROM chromedp/headless-shell

WORKDIR /data

COPY mdout /usr/local/bin

RUN mdout install theme -u https://github.com/JabinGP/mdout-theme-github/archive/0.1.1.zip -n github -k

ENTRYPOINT ["mdout"]