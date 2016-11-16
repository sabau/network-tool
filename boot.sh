#!/bin/bash

sleep 1
chvt 2

backtitle='QuiVIDEO Configurator'
exitstatus=0

restart_services(){
    /usr/bin/sudo systemctl restart nginx.service &
    /usr/bin/sudo systemctl restart php56-php-fpm.service &
    {
        i=0
        for ((i = 0 ; i <= 100 ; i+=20)); do
            sleep 1
            echo ${i}
        done
    } | whiptail --title "Restarting services" --backtitle "$backtitle" --gauge "Please wait while services restart" 6 60 0
}

create_portal(){
    name=$(whiptail --title "Portal Setup " --backtitle "$backtitle" --inputbox "Readable portal Name" 30 60 'Portal Name' 3>&1 1>&2 2>&3)
    exit=$?
    if [ ! ${exit} = 0 ]; then
        return
    fi
    url=$(whiptail --title "Portal Setup" --backtitle "$backtitle" --inputbox "Insert the portal URL (same as license if you have a portal-specifc license)" 30 60 'portal.quivideo.it' 3>&1 1>&2 2>&3)
    exit=$?
    if [ ! ${exit} = 0 ]; then
        return
    fi
    domain=$(whiptail --title "Portal Setup" --backtitle "$backtitle" --inputbox "Domain extension" 30 60 '.videoconferenza-hd2.com' 3>&1 1>&2 2>&3)
    exit=$?
    if [ ! ${exit} = 0 ]; then
        return
    fi
    port=$(whiptail --title "Portal Setup " --backtitle "$backtitle" --radiolist "Select port and schema" 30 60 2 "80" "http" ON "443" "https" OFF 3>&1 1>&2 2>&3)
    exit=$?
    if [ ! ${exit} = 0 ]; then
        return
    fi
    user=$(whiptail --title "Super Credentials" --backtitle "$backtitle" --inputbox "Insert superuser username" 30 60 'super' 3>&1 1>&2 2>&3)
    exit=$?
    if [ ! ${exit} = 0 ]; then
        return
    fi
    pass=$(whiptail --title "Super Credentials" --backtitle "$backtitle" --inputbox "Insert superuser password" 30 60 '*' 3>&1 1>&2 2>&3)
    exit=$?
    if [ ! ${exit} = 0 ]; then
        return
    fi

    #Create portal and sync tenants
    /usr/bin/sudo -u nginx /bin/bash -c "/var/www/applications/dashboard/current/bin/app-incubator portal-create '$name' $url $port $domain $user $pass"
    /usr/bin/sudo -u nginx /bin/bash -c "/var/www/applications/dashboard/current/bin/app-incubator portal-tenant-autosync "
}

insert_license() {
    old_license=$(</var/www/applications/dashboard/incubator_dashboard)
    #echo ${old_license}
    /usr/bin/sudo -u nginx /bin/bash -c "touch /var/www/applications/dashboard/incubator_dashboard.new"

    pass=$(/usr/bin/sudo -u nginx dialog --title "Insert License" --backtitle "$backtitle" --editbox /var/www/applications/dashboard/incubator_dashboard.new 30 60 3>&1 1>&2 2>&3)
    exit=$?
    if [ ! ${exit} = 0 ]; then
        return
    fi

    whiptail --title "Insert License" --backtitle "$backtitle" --yesno "Remove the old license and save the new one?" 30 60 --no-button "Undo" --yes-button "Save"
    ret=$?
    if [ ${ret} = 0 ]; then
        /usr/bin/sudo -u nginx /bin/bash -c "rm -rf /var/www/applications/dashboard/incubator_dashboard.old"
        /usr/bin/sudo -u nginx /bin/bash -c "cp /var/www/applications/dashboard/incubator_dashboard /var/www/applications/dashboard/incubator_dashboard.old"
        #/usr/bin/sudo -u nginx /bin/bash -c "mv /var/www/applications/dashboard/incubator_dashboard.new /var/www/applications/dashboard/incubator_dashboard"
        /usr/bin/sudo -u nginx /bin/bash -c "cat > /var/www/applications/dashboard/incubator_dashboard << EOL
${pass}
EOL"
    fi
}

fetch_data() {
    start=$(dialog --title "Starting date" --backtitle "$backtitle" --date-format '%Y%m%d' --calendar "Please choose the date where we start gather data" 30 60 `date -d '-1 month' +"%d %m %Y"` 3>&1 1>&2 2>&3)
    echo "Start: ${start}"
    exit=$?
    if [ ! ${exit} = 0 ]; then
        return
    fi
    if [ $(date -u +'%s') -le $(date -ud ${start} +'%s') ]; then
        return
    fi

    end=$(dialog --title "Ending date" --backtitle "$backtitle" --date-format '%Y%m%d' --calendar "Please choose the date where we finish gather data" 30 60 `date -d '-1 day' +"%d %m %Y"` 3>&1 1>&2 2>&3)
    echo "End: ${end}"
    exit=$?
    if [ ! ${exit} = 0 ]; then
        return
    fi
    if [ $(date -u +'%s') -le $(date -ud ${end} +'%s') ]; then
        return
    fi


    if [ $(date -ud ${end} +'%s') -le $(date -ud ${start} +'%s') ]; then
        tmp=${start}
        start=${end}
        end=${tmp}
    fi
    d1=$(date -ud ${end} +'%s')
    d2=$(date -ud ${start} +'%s')
    length=$(( ( d1 - d2 )/86400 ))
    start=$(date -ud ${start} +'%Y-%m-%d')
    end=$(date -ud ${end} +'%Y-%m-%d')

    echo "Start on ${start} until ${end}, we are fetching ${length} days of CDRs"
    /usr/bin/sudo -u nginx /bin/bash -c "/var/www/applications/dashboard/current/bin/app-incubator cdr 1 $start $end" &
    /usr/bin/sudo -u nginx /bin/bash -c "/var/www/applications/dashboard/current/bin/app-incubator fetch-interval $start $end" &
    {
        i=0
        for ((i = 0 ; i <= $length ; i+=1)); do
            sleep 1
            percentage=$((200*$i/$length % 2 + 100*$i/$length))
            echo ${percentage}
        done
    } | whiptail --title "Fetching ${length} days" --backtitle "$backtitle" --gauge "Please wait while all data is gathered" 6 60 0
}

set_net() {
    #######################################
    # live change
    # plus static change that persist
    #######################################

    /usr/bin/sudo sed -i.bak 's/BOOTPROTO=dhcp/BOOTPROTO=static/' /etc/sysconfig/network-scripts/ifcfg-eth0

    #update Hostname
    hostname=$(hostname)
    newhostname=$(whiptail --title "Set Hostname" --backtitle "$backtitle" --inputbox "Set your FQDN address" 30 60 ${hostname} 3>&1 1>&2 2>&3)
     exit=$?
    if [ ! $exit = 0 ]; then
        echo "Canceled NEW HOSTNAME action"
        return
    fi
    echo "New hostname will be " $newhostname
    /usr/bin/sudo hostnamectl set-hostname $newhostname
    /usr/bin/sudo sed -i "s*$hostname*$newhostname*" /etc/hostname

    declare -A comp
    comp[IPADDR]="Set your IP Address"
    comp[GATEWAY]="Set your Gateway Address"
    comp[PREFIX]="Set Netmask Prefix"
    comp[DNS1]="Set DNS1 Address"
    comp[DNS2]="Set DNS2 Address"

    for i in "${!comp[@]}"
    do
      echo "key: $i: ${comp[$i]}"

      ipaddr=$(grep -i $i /etc/sysconfig/network-scripts/ifcfg-eth0|awk -F= '{print $2}')
      newip=$(whiptail --title "${comp[$i]}" --backtitle "$backtitle" --inputbox "${comp[$i]}" 30 60 ${ipaddr} 3>&1 1>&2 2>&3)
      exit=$?
      if [ ! $exit = 0 ]; then
          echo "Canceled $i action"
          /usr/bin/sudo rm -rf /etc/sysconfig/network-scripts/ifcfg-eth0
          /usr/bin/sudo mv /etc/sysconfig/network-scripts/ifcfg-eth0.bak /etc/sysconfig/network-scripts/ifcfg-eth0
          return
      fi
      echo "New $i will be " $newip
      /usr/bin/sudo sed -i "s|^$ipaddr |^$newip |" /etc/hosts
      if grep -Fq "$i" /etc/sysconfig/network-scripts/ifcfg-eth0
      then
          /usr/bin/sudo sed -i "s/$i.$ipaddr/$i=$newip/" /etc/sysconfig/network-scripts/ifcfg-eth0
      else
          /usr/bin/sudo echo "$i=$newip" >> /etc/sysconfig/network-scripts/ifcfg-eth0
      fi

    done

    #/etc/init.d/network restart
    /usr/bin/sudo systemctl restart network
}



while [ ${exitstatus} = 0 ]
do

OPTION=$(whiptail --cancel-button "Quit" --title "QuiVIDEO Main Menu" --backtitle "$backtitle" --menu "Select which action you intend to do" 30 60 7 "1" "Set machine IP" "2" "Restart web services" "3" "Set dashboard server name" "4" "Update Dashboard license" "5" "Add Portal to dashboard" "6" "Show Dashboard License" "7" "Fetch Data"  3>&1 1>&2 2>&3)

exitstatus=$?
if [ ${exitstatus} = 0 ]; then
    case ${OPTION} in
        1)
        set_net
        ;;
        2)
        restart_services
        ;;
        3)
        arr=($(awk '/^\s*server_name/ {f=1} f {p=$NF;sub(/;$/,"",p);print p} /;$/ {f=0}' /etc/nginx/conf.d/dashboard.conf))
        hostname=${arr[0]}
        newhostname=$(whiptail --title "Set dashboard URL" --backtitle "$backtitle" --inputbox "Set your desired address" 30 60 ${hostname} 3>&1 1>&2 2>&3)
        exit=$?
        if [ $exit = 0 ]; then
            echo "New dashboard URL will be " ${newhostname}
            sudo sed -i.bak "s*$hostname*$newhostname*g" /etc/hosts
            sudo sed -i.bak "s*$hostname*$newhostname*g" /etc/nginx/conf.d/dashboard.conf
            sudo sed -i.bak "s*$hostname*$newhostname*g" /var/www/applications/dashboard/current/config/autoload_incubator/global.php
            restart_services
        else
            echo "Canceled CHANGE URL action"
        fi
        ;;
        4)
        insert_license
        ;;
        5)
        create_portal
        ;;
        6)
        pass=$(whiptail --title "Show License" --backtitle "$backtitle" --textbox /var/www/applications/dashboard/incubator_dashboard 30 60 --scrolltext 3>&1 1>&2 2>&3)
        ;;
        7)
        fetch_data
        ;;
        *)
        echo "Option not yet implemented" $OPTION
    esac
else
    whiptail --title "Quit Configurator"  --backtitle "$backtitle" --msgbox "Bye! Thanks for choosing QuiVIDEO. Hit OK to logout." 30 60
    clear
    exit
fi
done
clear
exit

