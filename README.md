= nds =

NftablesDynamicSets

Monitor Nftable Sets defined under
`/etc/nft.conf.d/sets.d/domains.d/*.conf`,
treat the set name as an FQDN, look it up,
update the set with the resulting IPs, write the set back to disk and
reload nftables atomically with `nft -f /etc/nft.conf`.

== Assumptions ==

- Set files for domains are stored as
    /etc/nft.conf.d/sets.d/domains.d/*.conf
- Set files follow the example syntax:
    table ip filter {
      set google.ca {
        type ipv4_addr
        elements = { 8.8.8.8 }
      }
    }
