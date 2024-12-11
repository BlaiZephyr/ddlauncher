{
  lib,
  buildGoModule,
  xorg,
  pkg-config,
  libGL,
}:
buildGoModule {
  pname = "ddlauncher";
  version = "0.0.1";

  src = ./.;

  nativeBuildInputs = [ pkg-config ];

  buildInputs = [
    xorg.libX11
    xorg.libXcursor
    xorg.libXrandr
    xorg.libXinerama
    xorg.libXi
    xorg.libXxf86vm
    libGL
  ];

  vendorHash = "sha256-iK9JMDkRGGqKeUSoiQqwFDqOrUaGlpWQdIQg888vJ4g=";

  meta = {
    description = "DDNet Client Manager";
    homepage = "https://github.com/BlaiZephyr/ddlauncher";
    mainProgram = "ddlauncher";
    platforms = lib.platforms.linux;
    maintainers = with lib.maintainers; [ theobori ];
  };
}
