#!/usr/bin/env bash
# Determine the OS and architecture

# Set up a cleanup function to be triggered upon script exit
__cleanup ()
{
   rm "${p_name}.tar.gz" 2>/dev/null
   if [ -n "$file" ]; then
       rm $file 2>/dev/null
   fi
   # base=$(basename $(pwd))
   # parent=$(dirname $(pwd))
   # if [ "${base}" != "installer" ] || [ "${parent}" != "aion-cli" ]; then
   #     rm install.sh 2>/dev/null
   # fi
}

trap __cleanup EXIT

OS=$(uname -s)
ARCH=$(uname -m)

file=""
p_name="aion-cli"

url="https://github.com/Slug-Boi/${p_name}/releases/latest/download/"

TAG=$(curl -L -s https://api.github.com/repos/Slug-Boi/${p_name}/releases/latest | grep -m1 '"v.*' | cut -c 16- | rev | cut -c 3- | rev)

# Set the download URL based on the OS and architecture
if [ "$OS" == "Linux" ]; then
  URL="${url}${p_name}-${TAG}-linux.tar.gz"
    file="$p_name"
elif [ "$OS" == "Darwin" ]; then
    if [ "$ARCH" == "x86_64" ]; then
        URL="${url}${p_name}-${TAG}-macos-x86_64.tar.gz"
        file="$p_name"
    else
        URL="${url}${p_name}-${TAG}-macos-aarch64.tar.gz"
        echo ${TAG}
        echo ${URL}
        file="$p_name"
    fi
else
    echo "Unsupported OS: $OS"
    exit 1
fi

# Download and run the script
echo ${p_name}
curl -L -o ${p_name}.tar.gz $URL && \
  tar -xvzf ${p_name}.tar.gz && \
rm ${p_name}.tar.gz && \
chmod +x $file && ./$file -h
if [ $? -ne 0 ]; then
    echo "Failed to extract the binary"
    exit 1
fi

# Move the binary to the current directory
read -p "Enter the directory to move the binary to (default: /usr/local/bin/${p_name}): " target_dir
target_dir=${target_dir:-/usr/local/bin/$p_name}

if [ ! -d "$(dirname "$target_dir")" ]; then
    echo "Directory does not exist: $(dirname "$target_dir")"
    exit 1
fi

sudo mv $file "$target_dir"
if [ $? -ne 0 ]; then
    echo "Failed to move the binary to $target_dir"
    exit 1
fi

echo "Binary moved to $target_dir successfully"
