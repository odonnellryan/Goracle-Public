useradd $USERNAME
echo $PASSWORD | passwd $USERNAME --stdin
