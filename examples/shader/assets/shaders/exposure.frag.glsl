#version 330 core

// Keep in mind, I have very little idea what I'm doing when it comes
// to these shaders, so take what you see here with a grain of salt.

out vec4 fragColor;

// Pixel default uniforms
uniform vec4 uTexBounds;
uniform sampler2D uTexture;
uniform sampler2D uBackBuffer;

// Our custom uniforms
uniform float uAmount;

void main() {
    // It is often very useful to normalize the fragment coordinate. Usually
    // represented as "uv" we do so here:
    vec2 uv = gl_FragCoord.xy / uTexBounds.zw;
    fragColor = texture(uTexture, uv);
 
    // uAmount is programmed to be adjustable with the left and right keys
    // inside of Pixel
    fragColor *= texture(uBackBuffer, uv).a * uAmount;
}
