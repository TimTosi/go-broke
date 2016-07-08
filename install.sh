#! /bin/bash

# Install libsodium
wget https://download.libsodium.org/libsodium/releases/libsodium-1.0.3.tar.gz && \
wget https://download.libsodium.org/libsodium/releases/libsodium-1.0.3.tar.gz.sig && \
wget https://download.libsodium.org/jedi.gpg.asc && \
tar zxvf libsodium-1.0.3.tar.gz && \
cd libsodium-1.0.3 && \
./configure; make check && \
sudo make install && \
sudo ldconfig

# Install zeromq
git clone git@github.com:zeromq/libzmq.git && \
cd libzmq && \
sudo apt-get install autoconf && \
sudo apt-get install libtool && \
./autogen.sh && \
./configure --with-libsodium && \
make check && \
sudo make install && \
sudo ldconfig 
