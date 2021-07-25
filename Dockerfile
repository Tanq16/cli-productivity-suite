FROM ubuntu:20.04

# Set environment variables
ENV RUNNING_IN_DOCKER true
ENV SHELL /bin/zsh
ENV TERM xterm

# Enable multiverse repository
RUN sed -i "/^# deb.*multiverse/ s/^# //" /etc/apt/sources.list

# Update and Upgrade
RUN add-apt-repository ppa:longsleep/golang-backports; \
	apt update && apt upgrade -y

# Install other packages
RUN DEBIAN_FRONTEND="noninteractive" \
    apt install -y --no-install-recommends \
    build-essential \
    zlib1g-dev \
    libncurses5-dev \
    libgdbm-dev \
    libnss3-dev \
    libssl-dev \
    libreadline-dev \
    libffi-dev \
    apt-transport-https \
    software-properties-common \
    vim \
    curl \
    strace \
    ltrace \
    bat \
    fd-find \
    wget \
    gdb \
    git \
    default-jre default-jdk \
    python3-pkg-resources python3-setuptools python3-pip python3 python3-dev \
    php \
    nasm \
    ncat \
    nmap \
    openssl \
    tmux \
    tree \
    unzip \
    fzf \
    iproute2 \
    openssh-server

# Install zsh and tmux
RUN apt install -y --no-install-recommends zsh
RUN curl -L http://install.ohmyz.sh | sh
RUN git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k 2>/dev/null
RUN sed -i "s/robbyrussell/powerlevel10k\/powerlevel10k/" ~/.zshrc
RUN git clone https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
RUN git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null
RUN sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
RUN git clone https://github.com/seebi/dircolors-solarized.git dirco_for_script_color_option 2>/dev/null
RUN cp dirco_for_script_color_option/dircolors.256dark ~/.oh-my-zsh/dircolors.256dark
RUN rm -rf dirco_for_script_color_option/
RUN git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 2>/dev/null
RUN git clone https://github.com/jimeh/tmux-themepack.git ~/.tmux-themepack
RUN wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 2>/dev/null
RUN mv tmuxconf ~/.tmux.conf
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 2>/dev/null
RUN wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/add_to_rc 2>/dev/null
RUN cat add_to_rc >> ~/.zshrc
RUN rm add_to_rc
RUN cp ~/.zshrc temptemp
RUN cat temptemp | grep -vE "^#" | grep -vE "^$" > ~/.zshrc
RUN chsh -s /usr/bin/zsh
RUN echo "[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh" >> ~/.zshrc

# Install vim extensions
RUN git clone https://github.com/itchyny/lightline.vim ~/.vim/pack/plugins/start/lightline 2>/dev/null
RUN git clone https://github.com/jiangmiao/auto-pairs.git ~/auto_pairs_vim 2>/dev/null
RUN mkdir -p ~/.vim/plugin
RUN cp ~/auto_pairs_vim/plugin/auto-pairs.vim ~/.vim/plugin/
RUN rm -rf ~/auto_pairs_vim
RUN wget https://www.vim.org/scripts/download_script.php\?src_id\=21752 2>/dev/null
RUN mv 'download_script.php?src_id=21752' ~/.vim/supertab.vmb
RUN vim -c 'so %' -c 'q' ~/.vim/supertab.vmb
RUN wget https://raw.githubusercontent.com/dylnmc/novum.vim/master/colors/novum.vim 2>/dev/null
RUN mkdir -p ~/.vim/colors
RUN mv novum.vim ~/.vim/colors/novum.vim
RUN curl -fLo ~/.vim/plugin/NERD_Commenter.vim --create-dirs https://raw.githubusercontent.com/preservim/nerdcommenter/master/plugin/NERD_commenter.vim 2>/dev/null
RUN curl -fLo ~/.vim/doc/NERD_Commenter.txt --create-dirs https://raw.githubusercontent.com/preservim/nerdcommenter/master/doc/NERD_commenter.txt 2>/dev/null
RUN wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/.vimrcfile 2>/dev/null
RUN mv .vimrcfile ~/.vimrc
RUN sleep 2

# Write stuff to do into a file
RUN echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
RUN echo 'root:docker' | chpasswd 
COPY ./p10k.zsh .
RUN mv /p10k.zsh ~/.p10k.zsh
