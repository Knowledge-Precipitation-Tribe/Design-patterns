package FactoryMethodPattern;

public class VolunteerFactory implements IFactory {
    @Override
    public LeiFeng createLeifeng() {
        return new Volunter();
    }
}
