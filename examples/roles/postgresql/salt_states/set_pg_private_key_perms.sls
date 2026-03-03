{% import_yaml slspath|path_join('../package.yaml') as role_metadata %}
{% set role_type_id = role_metadata.package.TypeId %}
{% set roleInstance = pillar['local']|selectattr('RoleTypeId', '==', role_type_id )|first %}

postgres_SSLKey_permissions:
  cmd.run:
    - name: find {{roleInstance.Params.CertificatesDir}} -maxdepth 1 -name "*.key" -and -not - name "otel.key" -exec chmod 600 {}  \; | true
    - require:
      - Install role packages
